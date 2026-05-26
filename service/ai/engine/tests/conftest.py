import os
from pathlib import Path


os.environ.setdefault("APP_CONFIG", str(Path(__file__).with_name("test_config.yaml")))
