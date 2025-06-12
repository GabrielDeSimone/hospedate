Email-sender
==============

Este script de python permite enviar correos electrónicos utilizando un servidor SMTP. Está diseñado para ser invocado desde la línea de comandos o a través de otro programa.


Funcionalidad
-------------
Conecta con un servidor SMTP preconfigurado.
Envía un correo electrónico con un asunto, remitente, destinatario y cuerpo específico.


Requisitos
----------
Python 3.x - no requiere librerías adicionales de las standard
Acceso a un servidor SMTP y sus credenciales.


Configuración
-------------
Antes de utilizar el script, es necesario configurar las siguientes variables de entorno:
SMTP_SERVER: El servidor SMTP a utilizar (por ejemplo, "smtp.gmail.com").
SMTP_PORT: El puerto para el servidor SMTP (por ejemplo, 587 para STARTTLS).
SMTP_USERNAME: La dirección de correo electrónico o usuario para autenticarse en el servidor SMTP.
SMTP_PASSWORD: La contraseña para autenticarse en el servidor SMTP.


Uso
---
Para usar el script desde la terminal, navega al directorio `email-sender/` y ejecuta el siguiente comando:

```
python main.py <asunto> <remitente> <destinatario> <cuerpo>
```
Por ejemplo:

```
python main.py "Asunto del Correo" "mi-correo@example.com" "destinatario@example.com" "Este es el cuerpo del correo."
```

Outputs Esperados
-----------------
Éxito: Si el correo se envía correctamente, el script imprimirá Email sent successfully!.
Error: Si ocurre algún problema durante el envío, el script imprimirá un mensaje de error indicando la naturaleza del problema, por ejemplo: Error sending email: [descripción del error].