# Welcome!

Доступно 2 пользователя для проверки: 
- Login: ginkt95@gmail.com Pass: 1234 Role: admin
- Login: john95@gmail.com Pass: 1234 Role: user

Доступ к ресурсам осуществляется через адресную строку

`http://localhost:8000/{foo, bar, sigma}`

Панель авторизации находится по адресу

`http://localhost:8000/login`

# RBAC
Реализован в виде Middleware(функция RBACMiddleware).
Загружает роль текущего пользователя и проверяет наличие прав доступа к ресурсу. 


# Пример использования
![image](https://i.imgur.com/r0JMYFM.jpg)

Работа логгера:

![imagelogger](https://i.imgur.com/cL3UXuh.jpg)

# Docker
Для проверки приложения можно также использовать докер.
Сборка и запуск контейнера:

`docker build -t rbacapp .
docker run --publish 8000:8000 --name test --rm rbacapp`


