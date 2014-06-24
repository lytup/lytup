#!/usr/bin/env python

import os
import shutil
from datetime import datetime
from pymongo import MongoClient

UPLOAD_DIR = '/tmp'

client = MongoClient('mongodb://localhost')
db = client.lytup

# Delete from database
folder = db.folders.find_and_modify(
  query={'expiresAt': {'$lt': datetime.utcnow()}},
  remove=True
)

# Delete from file system
if folder is not None:
  print 'Delete folder {0}'.format(folder['id']) 
  try:
    shutil.rmtree(os.path.join(UPLOAD_DIR, folder['id']))
  except OSError, err:
    print err
