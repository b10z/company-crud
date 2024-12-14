DO $$
    DECLARE
        table_to_drop text;
    BEGIN
        FOR table_to_drop IN (SELECT table_name FROM information_schema.tables WHERE table_schema = 'xm_assessment' AND table_type = 'BASE TABLE')
            LOOP
                BEGIN
                    EXECUTE 'DROP TABLE ' || 'xm_assessment.' || table_to_drop || ' CASCADE';
                EXCEPTION
                    WHEN others THEN
                        RAISE NOTICE 'Error dropping table %', table_to_drop;
                END;
            END LOOP;
    END $$;