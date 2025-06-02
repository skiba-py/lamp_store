-- Обновляем пути к картинкам, убирая /static/ из пути
UPDATE products SET image = REPLACE(image, '/static/images/', '/images/') WHERE image LIKE '/static/images/%'; 