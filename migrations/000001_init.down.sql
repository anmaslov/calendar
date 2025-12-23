-- Drop indexes
DROP INDEX IF EXISTS idx_event_categories_event_id;
DROP INDEX IF EXISTS idx_event_attendees_event_id;
DROP INDEX IF EXISTS idx_events_status;
DROP INDEX IF EXISTS idx_events_exchange_id;
DROP INDEX IF EXISTS idx_events_end_time;
DROP INDEX IF EXISTS idx_events_start_time;

-- Drop tables
DROP TABLE IF EXISTS event_categories;
DROP TABLE IF EXISTS event_attendees;
DROP TABLE IF EXISTS events;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";

