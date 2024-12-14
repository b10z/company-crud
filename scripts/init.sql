CREATE SCHEMA IF NOT EXISTS xm_assessment;

CREATE TYPE xm_assessment.COMP_TYPE AS ENUM ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship');

CREATE TABLE xm_assessment.companies
(
    id               UUID        DEFAULT gen_random_UUID() PRIMARY KEY,
    name             VARCHAR(15)               NOT NULL UNIQUE,
    description      VARCHAR(3000),
    employees_number INT                       NOT NULL,
    is_registered    BOOLEAN                   NOT NULL,
    type             xm_assessment.COMP_TYPE                 NOT NULL,
    created_at       TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at       TIMESTAMPTZ               NOT NULL
);
