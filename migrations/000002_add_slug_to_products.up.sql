ALTER TABLE products ADD COLUMN slug VARCHAR(255);

UPDATE products
SET
    slug = TRIM(
        BOTH '-'
        FROM LOWER(
                REGEXP_REPLACE(
                    name, '[^a-zA-Z0-9]+', '-', 'g'
                )
            )
    );

ALTER TABLE products ALTER COLUMN slug SET NOT NULL;

ALTER TABLE products ADD CONSTRAINT product_slug_key UNIQUE (slug);