from fastapi.testclient import TestClient

from app.main import app
from app.services.agent import ALLOWED_TOOLS


def test_wind_summary_empty_evidence_returns_scaffold():
    client = TestClient(app)
    response = client.post("/v1/wind/summary/timeseries", json={"farm_code": "FY", "tower_code": "3"})

    assert response.status_code == 200
    body = response.json()
    assert body["evidence_count"] == 0
    assert body["status"] == "insufficient_evidence"
    assert "不能生成异常判断" in body["summary"]


def test_wind_timeseries_summary_computes_numeric_metrics():
    client = TestClient(app)
    response = client.post(
        "/v1/wind/summary/timeseries",
        json={
            "farm_code": "FY",
            "tower_code": "3",
            "evidence": [
                {
                    "points": [
                        {"ts": "2026-05-25 10:00:00", "values": {"x": "1.2"}},
                        {"ts": "2026-05-25 10:01:00", "values": {"x": "1.8"}},
                    ]
                }
            ],
        },
    )

    assert response.status_code == 200
    body = response.json()
    assert body["status"] == "ok"
    assert body["metrics"]["fields"]["x"]["count"] == 2
    assert body["metrics"]["fields"]["x"]["max"] == 1.8


def test_wind_alarm_summary_aggregates_levels():
    client = TestClient(app)
    response = client.post(
        "/v1/wind/summary/alarm",
        json={
            "farm_code": "FY",
            "evidence": [
                {
                    "list": [
                        {"towerCode": "03", "alarmLevel": 4, "status": 0},
                        {"towerCode": "03", "alarmLevel": 2, "status": 1},
                    ]
                }
            ],
        },
    )

    assert response.status_code == 200
    body = response.json()
    assert body["status"] == "ok"
    assert body["metrics"]["risk"] == "critical"
    assert body["metrics"]["alarm_count"] == 2


def test_wind_alarm_summary_uses_go_aggregate_evidence():
    client = TestClient(app)
    response = client.post(
        "/v1/wind/summary/alarm",
        json={
            "farm_code": "FY",
            "evidence": [
                {
                    "source": "tdengine.alarm",
                    "total": 12,
                    "sampled": 3,
                    "truncated": True,
                    "granularity": "1h",
                    "by_level": {"4": 2, "2": 10},
                    "by_status": {"0": 12},
                    "by_alarm_code_top": [{"key": "1001", "count": 5}],
                    "by_tower_top": [{"key": "03", "count": 7}],
                    "by_time_bucket": [{"bucket_start": "2026-05-25 10:00:00", "count": 6}],
                    "peak_bucket": {"bucket_start": "2026-05-25 10:00:00", "count": 6},
                    "samples": [
                        {"ts": "2026-05-25 10:01:00", "alarm_level": 4, "alarm_code": 1001, "status": 0},
                        {"ts": "2026-05-25 10:02:00", "alarm_level": 2, "alarm_code": 1002, "status": 0},
                    ],
                }
            ],
        },
    )

    assert response.status_code == 200
    body = response.json()
    assert body["status"] == "ok"
    assert body["metrics"]["alarm_count"] == 12
    assert body["metrics"]["risk"] == "critical"
    assert body["metrics"]["status_counts"] == {"0": 12}
    assert body["metrics"]["returned_records"] == 3


def test_wind_timeseries_summary_ignores_alarm_samples():
    client = TestClient(app)
    response = client.post(
        "/v1/wind/summary/timeseries",
        json={
            "farm_code": "FY",
            "evidence": [
                {
                    "source": "tdengine.alarm",
                    "total": 2,
                    "sampled": 2,
                    "by_level": {"4": 1, "2": 1},
                    "samples": [
                        {"ts": "2026-05-25 10:01:00", "alarm_level": 4, "alarm_code": 1001, "status": 0},
                        {"ts": "2026-05-25 10:02:00", "alarm_level": 2, "alarm_code": 1002, "status": 0},
                    ],
                }
            ],
        },
    )

    assert response.status_code == 200
    body = response.json()
    assert body["status"] == "insufficient_evidence"
    assert body["metrics"]["point_count"] == 0
    assert "fields" not in body["metrics"]


def test_wind_report_and_ticket_return_drafts():
    client = TestClient(app)
    payload = {"farm_code": "FY", "tower_code": "7", "evidence": [{"alarmCode": 1001, "alarmLevel": 3}]}

    report = client.post("/v1/wind/reports/health/draft", json=payload)
    ticket = client.post("/v1/wind/tickets/draft", json=payload)

    assert report.status_code == 200
    assert ticket.status_code == 200
    assert report.json()["status"] == "draft"
    assert ticket.json()["status"] == "draft"


def test_agent_tool_whitelist_is_wind_domain():
    expected = {
        "search_maintenance_sop",
        "query_sensor_timeseries",
        "query_alarm_events",
        "get_turbine_metadata",
        "compare_sensor_trend",
        "generate_health_report",
        "create_maintenance_ticket_draft",
    }

    assert ALLOWED_TOOLS == expected
