from fabric.contrib.project import *
from fabric.operations import *
from fabric.context_managers import *

def sync_web():
  rsync_project('/var/www/lytup/web', '../web/public')

def build_server():
  with lcd('../server'):
    local('GOARCH=amd64 GOOS=linux go build -o build/server')

def deploy_server():
  build_server()
  put('../server/build/server', '/var/www/lytup/server')

def build_web():
  with lcd('../web'):
    local('GOARCH=amd64 GOOS=linux go build -o build/web')

def deploy_web():
  build_web()
  put('../web/build/web', '/var/www/lytup/web')

def sync_scripts():
  rsync_project('/var/www/lytup', '../scripts')
