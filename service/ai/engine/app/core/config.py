import os
from functools import lru_cache
from pathlib import Path
from typing import Literal

import yaml
from pydantic import BaseModel, ConfigDict


class Settings(BaseModel):

    # 表示这个 Settings 对象创建之后不允许修改。
    model_config = ConfigDict(frozen=True)

    app_name: str = "ai-copilot-engine"
    app_env: Literal["dev", "staging", "prod"] = "dev"
    mock_llm: bool = True
    llm_provider: str = "mock"
    llm_model: str = "deepseek-v4-pro"
    deepseek_chat_url: str = "https://api.deepseek.com/chat/completions"
    llm_api_key: str = ""
    llm_timeout_seconds: int = 60
    postgres_dsn: str = ""
    elasticsearch_url: str = ""
    kb_chunks_index: str = "kb_chunks_index"
    security_logs_index: str = "security_logs_index"
    embedding_provider: Literal["mock", "dashscope"] = "mock"
    embedding_model: str = "text-embedding-v4"
    dashscope_embeddings_url: str = "https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings"
    embedding_api_key: str = ""
    embedding_dimension: int = 1024
    embedding_batch_size: int = 10
    embedding_timeout_seconds: int = 30
    rerank_provider: Literal["mock", "dashscope"] = "mock"
    rerank_model: str = "qwen3-rerank"
    dashscope_rerank_url: str = "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank"
    rerank_api_key: str = ""
    rerank_top_n: int = 5
    rerank_timeout_seconds: int = 30
    agent_go_rpc_target: str = "127.0.0.1:8091"
    agent_go_rpc_timeout_seconds: int = 30
    agent_max_steps: int = 6
    agent_checkpoint_backend: Literal["memory", "postgres"] = "postgres"
    agent_checkpoint_dsn: str = ""


# Mapping from flat Settings field names to (yaml_section, yaml_key) pairs.
# Bridges the nested YAML structure and the flat Settings interface so that
# callers can continue using settings.elasticsearch_url etc. unchanged.
_FIELD_MAP: dict[str, tuple[str, str]] = {
    "app_name": ("app", "name"),
    "app_env": ("app", "env"),
    "mock_llm": ("llm", "mock"),
    "llm_provider": ("llm", "provider"),
    "llm_model": ("llm", "model"),
    "deepseek_chat_url": ("llm", "deepseek_chat_url"),
    "llm_api_key": ("llm", "api_key"),
    "llm_timeout_seconds": ("llm", "timeout_seconds"),
    "postgres_dsn": ("postgres", "dsn"),
    "elasticsearch_url": ("elasticsearch", "url"),
    "kb_chunks_index": ("elasticsearch", "kb_chunks_index"),
    "security_logs_index": ("elasticsearch", "security_logs_index"),
    "embedding_provider": ("embedding", "provider"),
    "embedding_model": ("embedding", "model"),
    "dashscope_embeddings_url": ("embedding", "dashscope_embeddings_url"),
    "embedding_api_key": ("embedding", "api_key"),
    "embedding_dimension": ("embedding", "dimension"),
    "embedding_batch_size": ("embedding", "batch_size"),
    "embedding_timeout_seconds": ("embedding", "timeout_seconds"),
    "rerank_provider": ("rerank", "provider"),
    "rerank_model": ("rerank", "model"),
    "dashscope_rerank_url": ("rerank", "dashscope_rerank_url"),
    "rerank_api_key": ("rerank", "api_key"),
    "rerank_top_n": ("rerank", "top_n"),
    "rerank_timeout_seconds": ("rerank", "timeout_seconds"),
    "agent_go_rpc_target": ("agent", "go_rpc_target"),
    "agent_go_rpc_timeout_seconds": ("agent", "go_rpc_timeout_seconds"),
    "agent_max_steps": ("agent", "max_steps"),
    "agent_checkpoint_backend": ("agent", "checkpoint_backend"),
    "agent_checkpoint_dsn": ("agent", "checkpoint_dsn"),
}


def _default_config_path() -> Path:
    """Resolve the default config.yaml path relative to the project root.

    The project root is three directories up from this file:
        app/core/config.py  ->  project_root/
    """
    return Path(__file__).resolve().parent.parent.parent / "config.yaml"


def _flatten_yaml(raw: dict) -> dict:
    """Convert nested YAML dict into a flat dict keyed by Settings field names.

    Missing sections or keys are silently skipped so that Settings defaults
    fill in the gaps.
    """
    flat: dict = {}
    for field_name, (section, key) in _FIELD_MAP.items():
        section_data = raw.get(section)
        if isinstance(section_data, dict) and key in section_data:
            flat[field_name] = section_data[key]
    return flat


def load_settings(config_path: Path | None = None) -> Settings:
    """Load and validate settings from a YAML config file.

    Args:
        config_path: Absolute or relative path to the YAML file.
                     Defaults to config.yaml in the project root,
                     or the path given by the APP_CONFIG env var.

    Returns:
        A validated Settings instance.

    Raises:
        FileNotFoundError: If the config file does not exist.
        pydantic.ValidationError: If config values fail validation.
    """
    if config_path is None:
        # APP_CONFIG env var is an infrastructure-level concern
        # (which file to read), not an application config override.
        env_path = os.getenv("APP_CONFIG")
        if env_path:
            config_path = Path(env_path)
        else:
            config_path = _default_config_path()

    if not config_path.is_file():
        raise FileNotFoundError(f"Config file not found: {config_path}")

    raw: dict = yaml.safe_load(config_path.read_text(encoding="utf-8")) or {}
    flat = _flatten_yaml(raw)
    if not flat.get("embedding_api_key"):
        flat["embedding_api_key"] = os.getenv("DASHSCOPE_API_KEY", "")
    if not flat.get("rerank_api_key"):
        flat["rerank_api_key"] = os.getenv("DASHSCOPE_API_KEY", "")
    # 用 flat 这个字典里的数据，创建并校验一个 Settings 对象。
    return Settings.model_validate(flat)


@lru_cache
def get_settings() -> Settings:
    """Cached settings accessor. Called throughout the application.

    Returns the same Settings instance for the lifetime of the process.
    For tests or multi-environment loading, use load_settings() directly.
    """
    return load_settings()
