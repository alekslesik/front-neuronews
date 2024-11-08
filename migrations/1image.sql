CREATE TABLE IF NOT EXISTS  `image` (
  `image_id` int NOT NULL AUTO_INCREMENT,
  `image_path` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `image_size` int DEFAULT NULL,
  `image_name` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `image_alt` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`image_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO neuronews.image (image_id,image_path,image_size,image_name,image_alt) VALUES
	 (1,'/static/upload/dvizhenie-ia-mi-sergei-furgal-priznali-ekstremistskim.png',67376,'dvizhenie-ia-mi-sergei-furgal-priznali-ekstremistskim','Движение «Я/мы Сергей Фургал» признали экстремистским'),
	 (2,'/static/upload/microsoft-nachnet-pomechat-slabie-pk.png',60297,'microsoft-nachnet-pomechat-slabie-pk','Microsoft начнет помечать слабые ПК'),
	 (3,'/static/upload/microsoft-nachnet-pomechat-slabie-pk.png',60297,'microsoft-nachnet-pomechat-slabie-pk','Microsoft начнет помечать слабые ПК'),
	 (4,'/static/upload/v-moskovskii-muzei-garazh-prishli-sotrudniki-fsb.png',57935,'v-moskovskii-muzei-garazh-prishli-sotrudniki-fsb','В московский музей «Гараж» пришли сотрудники ФСБ'),
	 (5,'/static/upload/v-moskovskii-muzei-garazh-prishli-sotrudniki-fsb.png',57935,'v-moskovskii-muzei-garazh-prishli-sotrudniki-fsb','В московский музей «Гараж» пришли сотрудники ФСБ'),
	 (6,'/static/upload/nabiullina-predupredila-o-zamedlenii-ekonomiki.png',51808,'nabiullina-predupredila-o-zamedlenii-ekonomiki','Набиуллина предупредила о замедлении экономики');
