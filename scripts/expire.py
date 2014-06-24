#!/usr/bin/env python

import os
import shutil
from datetime import datetime
from pymongo import MongoClient

UPLOAD_DIR = '/tmp'

client = MongoClient('mongodb://localhost')
db = client.lytup

qry = {
  'expiresAt': {'$lt': datetime.now()},
  'status': {'$ne': 'EXPIRED'}
}
folders = db.folders.find(qry)

# Delete expired folders
for fol in folders:
  fol_path = os.path.join(UPLOAD_DIR, fol['id'])
  print 'Deleting folder ' + fol_path
  try:
    shutil.rmtree(os.path.join(UPLOAD_DIR, fol_path))
    # Update folder status
    db.folders.update({'id': fol['id']}, {'$set': {'status': 'EXPIRED'}})
  except OSError, err:
    print err

# Delete from database
db.folders.remove(qry)
