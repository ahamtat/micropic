# Проект "Превьювер изображений" [![Go Report Card](https://goreportcard.com/badge/github.com/AcroManiac/micropic)](https://goreportcard.com/report/github.com/AcroManiac/micropic) [![CircleCI Build Status](https://circleci.com/gh/AcroManiac/micropic.svg?style=shield)](https://circleci.com/gh/AcroManiac/micropic) [![GPL Licensed](https://img.shields.io/badge/license-GPL-green.svg)](https://raw.githubusercontent.com/AcroManiac/micropic/master/LICENSE)
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


## Масштабирование микросервисов
Рассмотрим пример, когда к сервису обращается большое количество
клиентов (например, при нагрузочном тестировании).
Тогда возможно возникновение ситуации, в которой производительность
микросервисов будет недостаточной для обработки всех запросов.

Чтобы промоделировать такую ситуацию, выполните следующую команду:

    $ go run ./cmd/warmer/*.go -number 100 -workers 2 -proxy http://localhost:8080/fill

HTTP-прокси при обработке запроса от клиента устанавливает
таймаут в 5 секунд. Если микросервисы не успевают обработать
запрос, то прокси возвращает ошибку HTTP 500 Internal Server Error:

    ...
    {"level":"error","time":"2020-04-30T17:18:32.605+0300","message":"returned error response","code":500,"body":"dGltZW91dCBlbGFwc2VkIG9uIFJNUVJQQyByZXF1ZXN0IHNlbmRpbmc="}
    {"level":"error","time":"2020-04-30T17:18:32.629+0300","message":"returned error response","code":500,"body":"dGltZW91dCBlbGFwc2VkIG9uIFJNUVJQQyByZXF1ZXN0IHNlbmRpbmc="}
    {"level":"error","time":"2020-04-30T17:18:32.689+0300","message":"returned error response","code":500,"body":"dGltZW91dCBlbGFwc2VkIG9uIFJNUVJQQyByZXF1ZXN0IHNlbmRpbmc="}
    ...
    {"level":"info","time":"2020-04-30T17:18:32.689+0300","message":"Application working time is 13.562350365s"}

Теперь выполним масштабирование сервиса, увеличив количество
его компонентов. Для этого выполните команду:

    $ make scale-up
      >  Scaling previewers up
    Starting deployments_previewer_1 ... done
    Creating deployments_previewer_2 ... done
    Creating deployments_previewer_3 ... done
    Creating deployments_previewer_4 ... done
    Creating deployments_previewer_5 ... done
    
Проверьте список микросервисов. Теперь он содержит реплики:

    $ docker ps
    CONTAINER ID        IMAGE                       COMMAND                  CREATED             STATUS              PORTS                                                                                        NAMES
    ce0e0012a8a2        deployments_previewer       "dockerize -wait tcp…"   11 seconds ago      Up 9 seconds                                                                                                     deployments_previewer_2
    55d1f4df0131        deployments_previewer       "dockerize -wait tcp…"   11 seconds ago      Up 9 seconds                                                                                                     deployments_previewer_5
    7220b2600b80        deployments_previewer       "dockerize -wait tcp…"   11 seconds ago      Up 10 seconds                                                                                                    deployments_previewer_3
    2306fd4b3d17        deployments_previewer       "dockerize -wait tcp…"   11 seconds ago      Up 8 seconds                                                                                                     deployments_previewer_4
    e1de710ef1d7        deployments_previewer       "dockerize -wait tcp…"   2 hours ago         Up 2 hours                                                                                                       deployments_previewer_1
    23212432d1e8        flashspys/nginx-static      "nginx -g 'daemon of…"   2 hours ago         Up 2 hours          0.0.0.0:80->80/tcp                                                                           deployments_imagesource_1
    3119713b0af2        deployments_cache           "/bin/cache"             2 hours ago         Up 2 hours          0.0.0.0:50051->50051/tcp                                                                     deployments_cache_1
    84e67e032fe3        rabbitmq:3.8.3-management   "docker-entrypoint.s…"   2 hours ago         Up 2 hours          4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, 15671/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp   deployments_rabbitmq_1
    464178e1679d        deployments_proxy           "dockerize -wait tcp…"   2 hours ago         Up 2 hours          0.0.0.0:8080->8080/tcp                                                                       deployments_proxy_1

Повторим запуск команды для нагрузочного тестирования.
Убедитесь, что логи больше не содержат ошибок и время
выполнения программы уменьшилось:

    $ go run ./cmd/warmer/*.go -number 100 -workers 2 -proxy http://localhost:8080/fill
    ...
    {"level":"debug","time":"2020-04-30T17:34:24.558+0300","message":"preview returned successfully"}
    {"level":"debug","time":"2020-04-30T17:34:24.706+0300","message":"preview returned successfully"}
    {"level":"debug","time":"2020-04-30T17:34:24.754+0300","message":"preview returned successfully"}
    {"level":"info","time":"2020-04-30T17:34:24.754+0300","message":"Application working time is 9.491443441s"}
