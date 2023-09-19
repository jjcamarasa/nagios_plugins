# Plugins para nagios
Plugins creados o modificados para su uso con nagios/naemon.
- Instalación de los binarios: 
Mover a un path personalizado de nagios o a la carpeta de nagios estandar del sistema (/usr/lib/nagios/plugins, /usr/local/nagios/libexec, etc.) y dar permisos de ejecución

## Plugins disponibles
1. check-load-percent (Golang)
Basado en *check-load* de https://github.com/atc0005/go-check-plugins.git
Devuelve el uso de cpu en porcentajes (como el nsclient) en base al numero de CPUs del sistema. 
Solo funciona en linux (lee los datos de /proc/loadavg). 
Devuelve *Performance Data* para graficas. Probado en Debian 10/11 y RedHat 5/6.
