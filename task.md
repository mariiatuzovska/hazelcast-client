# task

1. Установить и настроить Hazelcast [http://hazelcast.org/download/](http://hazelcast.org/download/).

2. Сконфигурировать и запустить 3 ноды (инстанса) объединенные в кластер [https://hazelcast.org/getting-started-with-hazelcast/](https://hazelcast.org/getting-started-with-hazelcast/).

3. Продемонстрируйте работу [Distributed Map](http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#map).

4. Устанавливая различные значения Backup проверьте сколько процентов данных из  Map теряется при одновременном отключении нескольких нод. Добейтесь отсутствия потерь данных [Backing Up Maps](http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#backing-up-maps).

```
<map name="default">
        <backup-count>0</backup-count>
        <async-backup-count>1</async-backup-count>
        <read-backup-data>true</read-backup-data>
 </map>
```

5. Продемонстрируйте работу Distributed Map with locks запустив примеры подсчета значений в цикле без блокировки; с пессимистической; и оптимистической блокировками [Locking Maps](http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#locking-maps).

6. Продемонстрируйте работу [Distributed Queue](http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#queue). C одной/нескольких нод идет запись, с других чтение.

7. Настройте [Bounded queue](http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#setting-a-bounded-queue). C одной/нескольких нод идет запись, с других чтение - проверьте что добавление блокируется при заполнении очереди.

8. Продемонстрируйте работу [Distributed Topic](http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#topic). C одной/нескольких нод идет запись, с других чтение.

9.  Продемонстрируйте работу [Distributed Lock](http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#lock).
