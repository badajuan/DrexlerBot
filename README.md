# Bot de Twitter en Golang

## **Descripción**

### Este es un bot de Twitter bastante simple desarrollado en Golang que permite publicar tweets automáticamente desde la consola de comandos a una cuenta de desarrollador de Twitter. Actualmente funciona en un loop continuo que pide que ingreses el texto del tweet y lo publica en tu cuenta usando la API oficial de Twitter (andando al 07/2023).<br> Además guarda las respuestas del servidor en archivos JSON con la fecha y hora de la publicación.<br><br>
#### *Posteriores actualizaciones pueden incluir el uso de una pequeña base de datos para almacenar frases a tuitear*<br><br> 
---

# Documentación generada con ChatGPT 👇

## **Requisitos**

Antes de ejecutar el bot, asegúrate de tener los siguientes requisitos:

1. Golang instalado en tu sistema: [Descargar Golang](https://golang.org/dl/)

2. Claves de acceso de la API de Twitter:
   - Necesitarás crear una aplicación en la plataforma de desarrolladores de Twitter para obtener las claves de acceso.
   - Coloca las claves de acceso en un archivo de texto (claves.txt) con el siguiente formato:
   ```
   consumerKey=YOUR_CONSUMER_KEY
   consumerSecret=YOUR_CONSUMER_SECRET
   accessToken=YOUR_ACCESS_TOKEN
   accessSecret=YOUR_ACCESS_SECRET
   bearerToken=YOUR_BEARER_TOKEN
   ```
   Asegúrate de reemplazar `YOUR_CONSUMER_KEY`, `YOUR_CONSUMER_SECRET`, `YOUR_ACCESS_TOKEN`, `YOUR_ACCESS_SECRET` y `YOUR_BEARER_TOKEN` con tus propias claves.

## **Cómo usar el bot**

1. Clona este repositorio a tu sistema local:

```
git clone https://github.com/tu_usuario/bot-twitter-golang.git
cd bot-twitter-golang
```

2. Ejecuta el bot con el siguiente comando:

```
go run main.go
```

3. El bot te pedirá que ingreses el texto del tweet que deseas publicar.

4. Después de publicar el tweet, el bot guardará la respuesta del servidor en un archivo JSON en la carpeta "responses" con la fecha y hora de la publicación.

5. El bot te preguntará si deseas seguir twiteando. Si respondes "Y" o "y", podrás ingresar otro tweet y publicarlo; si respondes "N" o "n" (o cualquier otra cosa) el bot terminará la ejecución.

## **Disclaimer**

Este bot es solo una implementación básica y no incluye funcionalidades avanzadas como manejo de errores exhaustivo, monitoreo de respuestas del servidor, entre otros. Si planeas utilizar este bot para una aplicación más grande o para uso en producción, te recomendamos agregar funcionalidades adicionales y manejo de errores más robusto.

¡Diviértete usando tu bot de Twitter en Golang! 🐦