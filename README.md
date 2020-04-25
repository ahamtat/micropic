# Проект "Превьювер изображений" [![CircleCI Build Status](https://circleci.com/gh/AcroManiac/micropic.svg?style=shield)](https://circleci.com/gh/AcroManiac/micropic) [![GPL Licensed](https://img.shields.io/badge/license-GPL-green.svg)](https://raw.githubusercontent.com/AcroManiac/micropic/master/LICENSE.md)
Микросервис изготовления превью для изображений, хранящихся на внешних ресурсах. Проектная работа по курсу OTUS «Разработчик Golang».

Техническое задание на сервис [Превьювер изображений](https://github.com/OtusGolang/final_project/blob/master/03-image-previewer.md)

## Сборка проекта
После выгрузки проекта из репозитория выполните команду:

    $ make build

В результате чего будут собраны образы Docker для микросервисов. Для проверки образов выполните команду:

    $ docker images
    REPOSITORY              TAG                 IMAGE ID            CREATED             SIZE
    deployments_proxy       latest              353f62bd7ed8        6 seconds ago       40MB
    deployments_previewer   latest              64503cb66993        7 seconds ago       37MB
    deployments_cache       latest              d6c8cd04631f        28 seconds ago      26.5MB
    deployments_builder     latest              f674951a4062        43 seconds ago      1.23GB
    golang                  alpine              dda4232b2bd5        18 hours ago        370MB
    alpine                  latest              f70734b6a266        32 hours ago        5.61MB

## Порядок запуска микросервисов
Для запуска микросервисов выполните следующую команду:

    $ make run

Проверить работу микросервисов можно командой:

    $ docker ps
    CONTAINER ID        IMAGE                       COMMAND                  CREATED             STATUS              PORTS                                                                                        NAMES
    ed992d61da01        deployments_previewer       "dockerize -wait tcp…"   15 seconds ago      Up 3 seconds                                                                                                     deployments_previewer_1
    983f9f3f4c2c        flashspys/nginx-static      "nginx -g 'daemon of…"   16 seconds ago      Up 15 seconds       0.0.0.0:80->80/tcp                                                                           deployments_imagesource_1
    518bc8530204        rabbitmq:3.8.3-management   "docker-entrypoint.s…"   16 seconds ago      Up 15 seconds       4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, 15671/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp   deployments_rabbitmq_1
    e7537efad24a        deployments_proxy           "dockerize -wait tcp…"   16 seconds ago      Up 2 seconds        0.0.0.0:8080->8080/tcp                                                                       deployments_proxy_1
    b084683cc068        deployments_cache           "/bin/cache"             16 seconds ago      Up 14 seconds       0.0.0.0:50051->50051/tcp                                                                     deployments_cache_1
    
Чтобы запустить интеграционные тесты, выполните команду:

    $ make test
      >  Making integration tests
    ...
    Creating deployments_integration_tests_1 ... done
    ................... 19
    
    
    7 scenarios (7 passed)
    19 steps (19 passed)
    46.637408ms
    testing: warning: no tests to run
    PASS
    ok  	github.com/AcroManiac/micropic/test/integration	0.057s

Чтобы остановить микросервисы, выполните команду:
    
    $ make down

## Очистка проекта
Для очистки проекта от Docker-образов микросервисов выполните команду:

    $ make clean
