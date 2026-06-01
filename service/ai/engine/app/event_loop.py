"""Event-loop factories used by local ASGI launchers."""

from __future__ import annotations

import asyncio


def selector_event_loop_factory() -> asyncio.AbstractEventLoop:
    """Return the Selector loop required by psycopg async connections on Windows."""

    return asyncio.SelectorEventLoop()
