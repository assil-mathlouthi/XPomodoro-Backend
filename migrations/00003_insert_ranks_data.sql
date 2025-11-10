-- +goose Up
INSERT INTO ranks (name, min_xp) VALUES
('Wood I', 0), ('Wood II', 50), ('Wood III', 100), ('Wood IV', 150), ('Wood V', 210), ('Wood VI', 280), ('Wood VII', 360), ('Wood VIII', 450), ('Wood IX', 550);

INSERT INTO ranks (name, min_xp) VALUES
('Bronze I', 660), ('Bronze II', 780), ('Bronze III', 910), ('Bronze IV', 1050), ('Bronze V', 1200), ('Bronze VI', 1360), ('Bronze VII', 1530), ('Bronze VIII', 1710), ('Bronze IX', 1900);

INSERT INTO ranks (name, min_xp) VALUES
('Silver I', 2100), ('Silver II', 2310), ('Silver III', 2530), ('Silver IV', 2760), ('Silver V', 3000), ('Silver VI', 3250), ('Silver VII', 3510), ('Silver VIII', 3780), ('Silver IX', 4060);

INSERT INTO ranks (name, min_xp) VALUES
('Gold I', 4350), ('Gold II', 4650), ('Gold III', 4960), ('Gold IV', 5280), ('Gold V', 5610), ('Gold VI', 5950), ('Gold VII', 6300), ('Gold VIII', 6660), ('Gold IX', 7030);

INSERT INTO ranks (name, min_xp) VALUES
('Platinum I', 7410), ('Platinum II', 7800), ('Platinum III', 8200), ('Platinum IV', 8610), ('Platinum V', 9030), ('Platinum VI', 9460), ('Platinum VII', 9900), ('Platinum VIII', 10350), ('Platinum IX', 10810);

INSERT INTO ranks (name, min_xp) VALUES
('Sapphire I', 11280), ('Sapphire II', 11760), ('Sapphire III', 12250), ('Sapphire IV', 12750), ('Sapphire V', 13260), ('Sapphire VI', 13780), ('Sapphire VII', 14310), ('Sapphire VIII', 14850), ('Sapphire IX', 15400);

INSERT INTO ranks (name, min_xp) VALUES
('Amethyst I', 15960), ('Amethyst II', 16530), ('Amethyst III', 17110), ('Amethyst IV', 17700), ('Amethyst V', 18300), ('Amethyst VI', 18910), ('Amethyst VII', 19530), ('Amethyst VIII', 20160), ('Amethyst IX', 20800);

INSERT INTO ranks (name, min_xp) VALUES
('Master I', 21450), ('Master II', 22110), ('Master III', 22780), ('Master IV', 23460), ('Master V', 24150), ('Master VI', 24850), ('Master VII', 25560), ('Master VIII', 26280), ('Master IX', 27010);


-- +goose Down
DELETE FROM ranks;
