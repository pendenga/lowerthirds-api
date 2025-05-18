DROP TABLE HymnVerses;
DROP TABLE Hymns;

CREATE TABLE Hymns (
   id CHAR(36) NOT NULL,
   page INT NOT NULL,
   language CHAR(3) NOT NULL,
   name CHAR(100) NOT NULL,
   translation_id CHAR(36) NULL,
   deleted_dt DATETIME NULL,
   inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (id)
);
INSERT INTO Hymns (id, language, page, name) VALUES ('fd5905bb-35a4-4a2f-9e29-041f58f3d1a9', 'eng', 243, 'Let Us All Press On');
INSERT INTO Hymns (id, language, page, name) VALUES ('3549ebe2-b6cc-4433-a0dd-365ec4113d38', 'spa', 158, 'Trabajemos hoy en la obra');
INSERT INTO Hymns (id, language, page, name) VALUES ('bb125745-55eb-448c-b255-dac7ef6444cc', 'eng', 66, 'Rejoice, the Lord is King!');
INSERT INTO Hymns (id, language, page, name) VALUES ('339dee9c-e944-4eb1-bdb8-bf7e1b9c411f', 'eng', 3, 'Now Let Us Rejoice');
INSERT INTO Hymns (id, language, page, name) VALUES ('dbb6cabf-9466-46f2-9cfd-f0e06aa62869', 'eng', 6, 'Redeemer of Israel');
INSERT INTO Hymns (id, language, page, name) VALUES ('a4d02b9f-bef8-47db-8765-4e8cee76bb64', 'spa', 5, 'Redentor de Israel');
UPDATE Hymns SET translation_id = 'a4d02b9f-bef8-47db-8765-4e8cee76bb64' WHERE id = 'dbb6cabf-9466-46f2-9cfd-f0e06aa62869';
UPDATE Hymns SET translation_id = '3549ebe2-b6cc-4433-a0dd-365ec4113d38' WHERE id = 'fd5905bb-35a4-4a2f-9e29-041f58f3d1a9';

CREATE TABLE HymnVerses (
    hymn_id CHAR(36) NOT NULL,
    verse_number INT NOT NULL,
    verse_lines TEXT,
    optional TINYINT(1) NOT NULL DEFAULT 0,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (hymn_id, verse_number, deleted_dt)
);
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('dbb6cabf-9466-46f2-9cfd-f0e06aa62869', 1, 'Redeemer of Israel, our only delight,
On whom for a blessing we call.
Our Shadow by day and our pillar by night,
Our King, our Deliv’rer, our all');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('dbb6cabf-9466-46f2-9cfd-f0e06aa62869', 2, 'We know he is coming to gather his sheep
And lead them to Zion in love,
For why in the valley Of death should they weep
Or in the lone wilderness rove?');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('dbb6cabf-9466-46f2-9cfd-f0e06aa62869', 3, 'How long we have wandered as strangers in sin
And cried in the desert for thee!
Our foes have rejoiced When our sorrows they’ve seen,
But Israel will shortly be free.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('dbb6cabf-9466-46f2-9cfd-f0e06aa62869', 4, 'As children of Zion, good tidings for us.
The tokens already appear.
Fear not, and be just, For the kingdom is ours.
The hour of redemption is near.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines, optional)
VALUES ('dbb6cabf-9466-46f2-9cfd-f0e06aa62869', 5, 'Restore, my dear Savior, the light of thy face;
Thy soul-cheering comfort impart;
And let the sweet longing For thy holy place
Bring hope to my desolate heart.', 1);
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines, optional)
VALUES ('dbb6cabf-9466-46f2-9cfd-f0e06aa62869', 6, 'He looks! and ten thousands Of angels rejoice,
And myriads wait for his word;
He speaks! and eternity, Filled with his voice,
Re-echoes the praise of the Lord', 1);

INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('a4d02b9f-bef8-47db-8765-4e8cee76bb64', 1, 'Oh Dios de Israel, te rendimos loor
a ti, nuestro gran Redentor,
de día la sombra, de noche la luz,
del mundo eres Rey y Señor.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('a4d02b9f-bef8-47db-8765-4e8cee76bb64', 2, 'Sabemos que vienes tu grey a juntar,
la cual has de guiar a Sión.
En valle de muerte no nos dejarás,
ni en la vasta desolación.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('a4d02b9f-bef8-47db-8765-4e8cee76bb64', 3, 'Hemos errado mucho, clamando a ti,
extraños, en yermos del mal.
Los malos se gozan de nuestro pesar,
mas libre Israel quedará.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('a4d02b9f-bef8-47db-8765-4e8cee76bb64', 4, 'Nos regocijamos, oh hijos de Dios;
las señas presentes están.
Seamos valientes y fieles al Rey;
se vislumbra la gran redención.');

INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('fd5905bb-35a4-4a2f-9e29-041f58f3d1a9', 1, 'Let us all press on in the work of the Lord,
That when life is o’er we may gain a reward;
In the fight for right let us wield a sword,
The mighty sword of truth.
Fear not, though the enemy deride;
Courage, for the Lord is on our side.
We will heed not what the wicked may say,
But the Lord alone we will obey.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('fd5905bb-35a4-4a2f-9e29-041f58f3d1a9', 2, 'We will not retreat, though our numbers may be few
When compared with the opposite host in view;
But an unseen pow’r will aid me and you
In the glorious cause of truth.
Fear not, though the enemy deride;
Courage, for the Lord is on our side.
We will heed not what the wicked may say,
But the Lord alone we will obey.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('fd5905bb-35a4-4a2f-9e29-041f58f3d1a9', 3, 'If we do what’s right we have no need to fear,
For the Lord, our helper, will ever be near;
In the days of trial his Saints he will cheer,
And prosper the cause of truth.
Fear not, though the enemy deride;
Courage, for the Lord is on our side.
We will heed not what the wicked may say,
But the Lord alone we will obey.');

INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('3549ebe2-b6cc-4433-a0dd-365ec4113d38', 1, 'Trabajemos hoy en la obra del Señor,
y ganemos así un hogar celestial.
En la lucha cruel empuñemos, sin temor,
la espada de la verdad.
Firmes y valientes en la lid,
todo enemigo confundid.
Lucharemos a vencer el error;
seguiremos sólo al Señor.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('3549ebe2-b6cc-4433-a0dd-365ec4113d38', 2, 'Nuestras filas chicas jamás desmayarán,
a pesar de las huestes que contenderán,
y del cielo, Cristo poder nos dará
en defensa de la verdad.
Firmes y valientes en la lid,
todo enemigo confundid.
Lucharemos a vencer el error;
seguiremos sólo al Señor.');
INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines)
VALUES ('3549ebe2-b6cc-4433-a0dd-365ec4113d38', 3, 'Toda obra buena aleja el temor,
pues tenemos en Cristo un gran Defensor.
En las duras pruebas nos da el valor
de luchar por la verdad.
Firmes y valientes en la lid,
todo enemigo confundid.
Lucharemos a vencer el error;
seguiremos sólo al Señor.');

/*
SELECT * FROM HymnVerses;
SELECT * FROM Hymns;
*/