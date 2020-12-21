CREATE TABLE IF NOT EXISTS messages(
   ID integer NOT NULL,
   RU text NOT NULL,
   ENG text NOT NULL,
   CONSTRAINT messages_pkey PRIMARY KEY (ID)
);

INSERT INTO messages (ID, RU, ENG) VALUES (1, 'Привет кожаный мешок! Введи название мусора, который хочешь утилизовать или выбери из списка с помощью команды /getwastetypes.', 'Hello little man! If you want to save your planet I can help you to recycling waste! Type waste or use /getwastetypes command!');
INSERT INTO messages (ID, RU, ENG) VALUES (2, 'Ни одного пункта сдачи не найдено :(', 'Recycling points not found :(');
INSERT INTO messages (ID, RU, ENG) VALUES (3, 'Выберите тип отхода:', 'Choose waste type:');
INSERT INTO messages (ID, RU, ENG) VALUES (4, 'Для получения пунктов сдачи необходимо определить геолокацию. Нажмите на кнопку:', 'To get waste recycling points, you need to determine the location. Click on the button:');
INSERT INTO messages (ID, RU, ENG) VALUES (5, 'Введенный текст не распознался как вид отхода. Можете повторить запрос либо воспользоваться командой /getwastetypes  для вывода кнопок с видами отходов.', 'The entered text was not recognized as a waste type. You can repeat the request or use the /getwastetypes command to display buttons with waste types.');
INSERT INTO messages (ID, RU, ENG) VALUES (6, 'Отправить геолокацию', 'Send location');