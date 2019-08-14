# request-generator

```
mkdir -p ~/go/src/stash.tutu.ru/opscore-workshop-admin
cd ~/go/src/stash.tutu.ru/opscore-workshop-admin
git clone ssh://git@depot.tutu.ru:7999/opscore-workshop-admin/request-generator.git
cd request-generator
make init
make watcher
```


**Важно!**

Не забудьте добавить в свой проект в Bitbucket юзера ro-user с правами на чтение:

https://stash.tutu.ru/projects/opscore-workshop-admin/permissions

Без этого сборка в ci/cd работать не будет
