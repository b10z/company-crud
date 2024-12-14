INSERT INTO xm_assessment.companies (name, description, employees_number, is_registered, type, updated_at)
SELECT
        'Company #' || generate_series,
        'Sample Description',
        floor(random()*100 + 1),
        RANDOM()::INT::BOOLEAN,
        CAST ((ARRAY['Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship'])[FLOOR(random() * 4 + 1)] AS xm_assessment.COMP_TYPE),
        NOW() - (random() * interval '365 days')
FROM generate_series(1, 200);

