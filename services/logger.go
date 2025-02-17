package services

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.New(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,  // Taille max en MB avant rotation
		MaxBackups: 5,   // Nombre max de fichiers de log sauvegardés
		MaxAge:     30,  // Nombre de jours avant suppression des anciens logs
	}, "[APP] ", log.Ldate|log.Ltime|log.Lshortfile)
}