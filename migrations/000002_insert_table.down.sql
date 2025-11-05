-- 000002_insert_table.down.sql
TRUNCATE TABLE 
    payments,
    bookings,
    schedules,
    fields,
    venues,
    users,
    roles
RESTART IDENTITY CASCADE;