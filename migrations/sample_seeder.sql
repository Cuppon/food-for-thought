INSERT INTO source (description, location)
VALUES ('Tried & True', 'https://www.triedandtruerecipe.com/easy-miso-sesame-chickpeas/'),
       ('American', '/Users/ren/downloads/flag_usa.png');

INSERT INTO recipe_source (recipe_id, emoji_source_id)
VALUES (1, 1), (1, 2);

INSERT INTO ingredient (english_name, native_name, translated_name, shopping_link, english_category)
VALUES ('sesame oil', '', '',
        'https://www.amazon.com/Kadoya-Pure-Sesame-Fluid-Ounce/dp/B002HMN6SC/ref=sr_1_6?crid=28NMH7K348FR7&dib=eyJ2IjoiMSJ9.THeyXeGlGxLPjD70HFYA2twiZXkU-HhBTCinbkkRBudKHxKjOq9oLU62hgr0ebVNLUxfUWmHiqazeiUxfuALgDqH0U0K0yEKkceYsd_kzxN8lbDvktAz8E7ziwVo6JDa3EU8_ADFx9T6EFNNn66vCICj2sbce76b6xH2sSprVgO8TJwZOM5Swwc-2qSYaPyzuVk4eatIo9WQpHeKilZX6iCL0P9XzTYHyKkUEYDVGLRGPHrD7LynE6hzTJrCQFHiZTtA8qq07JwWeNfAtDPrKmbHi-_N98qK0cO3MeNKx84.3EeCj1j-rTGeCjJs70Ywz7fic4IzByNDp5Mh-OciC-0&dib_tag=se&keywords=sesame+oil&qid=1717517507&rdc=1&sprefix=sesame+oil%2Caps%2C188&sr=8-6', ''),
       ('shallot', '', '', '', ''),
       ('couscous (pearl)', '', '', '', 'couscous'),
       ('chickpeas', '', '', '', ''),
       ('quinoa', '', '', '', ''),
       ('miso (white/yuzu)', '', '', '', ''),
       ('water', '', '', '', ''),
       ('salt', '', '', '', ''),
       ('pepper', '', '', '', ''),
       ('shichimi togarashi', '七味唐辛子', 'shichi-mi tōgarashi', '', ''),
       ('scallion', '', '', '', '');

INSERT INTO recipe (attribution_source_id, cuisine_source_id, english_name, native_name, note, instruction)
VALUES (1, 2, 'Miso Sesame Chickpeas & Quinoa', '', null, ARRAY [
    '{
      "part":"main part",
      "steps":[
        {
          "is_parallel":false,
          "action":{
            "message":"Heat the sesame oil in a wide pot over medium heat. Once hot, add the shallot and cook for 3 minutes until it just begins to soften. Season with salt.",
            "unit_temperature":null
          },
          "note":null
        },
        {
          "is_parallel":false,
          "action":{
            "message":"Add the pearl couscous and cook for 3 minutes more until it begins to toast and turn golden brown.",
            "unit_temperature":null
          },
          "note":null
        },
        {
          "is_parallel":false,
          "action":{
            "message":"Add the chickpeas and season all over with salt and pepper.",
            "unit_temperature":null
          },
          "note":null
        },
        {
          "is_parallel":false,
          "action":{
            "message":"Add 1/3 of the water, let it heat up, then add miso paste and stir to dissolve",
            "unit_temperature":null
          },
          "note":null
        },
        {
          "is_parallel":true,
          "action":{
            "message":"In a separate small pot, cook the quinoa in water",
            "unit_temperature":null
          },
          "note":[
            {
              "message":"for high altitude, cook at",
              "unit_temperature":350
            }
          ]
        },
        {
          "is_parallel":false,
          "action":{
            "message":"Add the remaining 2/3 of the water, bring to boil, then simmer until done",
            "unit_temperature":null
          },
          "note":null
        },
        {
          "is_parallel":false,
          "action":{
            "message":"Mix in quinoa to couscous/chickpea pan, add the shichimi togarashi, additional sesame oil, and salt pepper to taste",
            "unit_temperature":null
          },
          "note":null
        },
        {
          "is_parallel":false,
          "action":{
            "message":"Finally, garnish with scallion",
            "unit_temperature":null
          },
          "note":null
        }
      ]
    }'::jsonb
    ]);

INSERT INTO ingredient_specification (component, recipe_id, ingredient_id, note, amount_quantity, amount_mass,
                                      preparation_quantity, preparation_type, preparation_length)
VALUES (null, 1, 1, 'not toasted', null, 2, 15, null, null),
       (null, 1, 2, 'large', null, null, 1, 'peeled & diced', null),
       (null, 1, 3, null, null, 1, 240, null, null),
       (null, 1, 4, null, null, 1, 425.23, 'drained & rinsed', null),
       (null, 1, 5, null, null, 1, 120, null, null),
       (null, 1, 6, 'white/yuzu', null, 2, 15, null, null),
       (null, 1, 7, 'more if needed', null, 2, null, null, null),
       (null, 1, 8, 'to taste', null, null, null, null, null),
       (null, 1, 9, 'to taste', null, null, null, null, null),
       (null, 1, 10, 'to taste', null, null, null, null, null),
       (null, 1, 11, 'for garnish', null, null, null, null, null);