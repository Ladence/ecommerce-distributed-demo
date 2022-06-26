CREATE TABLE IF NOT EXISTS events(
    id varchar primary key,
    eventtype varchar,
    aggregateid varchar,
    aggregatetype varchar,
    eventdata jsonb,
    stream varchar
)
