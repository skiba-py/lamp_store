-- Возвращаем /static/ в пути к картинкам
UPDATE products SET image = REPLACE(image, '/images/', '/static/images/') WHERE image LIKE '/images/%'; 