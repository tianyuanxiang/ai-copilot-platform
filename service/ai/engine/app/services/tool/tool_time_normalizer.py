from dataclasses import dataclass

@dataclass(frozen=True)
class NormalizedTimeRange:
    start_time: str
    end_time: str
    source_text: str
    label: str
    locked: bool = True

