BEGIN;

INSERT INTO characters (id, name) VALUES
(1, "ドラゴン"),
(2, "ナイト"),
(3, "マジシャン");

INSERT INTO gacha_probabilities (character_id, probability) VALUES
(1, 0.10),
(2, 0.20),
(3, 0.70);

COMMIT;
