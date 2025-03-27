DO $$ 
DECLARE 
    start_date DATE := CURRENT_DATE;
    end_date DATE := CURRENT_DATE + INTERVAL '1 year';
    d DATE;
BEGIN
    FOR d IN SELECT generate_series(start_date, end_date, '1 day') LOOP
        INSERT INTO available_slots (date, time, is_active)
        SELECT d, TIME '09:00:00' + INTERVAL '1 hour 30 minutes' * generate_series(0, 6), TRUE
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;
