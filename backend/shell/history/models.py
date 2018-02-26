from django.db import models

# Create your models here.


class Command(models.Model):
    hostname = models.CharField(max_length=255)
    timestamp = models.DateTimeField()
    username = models.CharField(max_length=32)
    altusername = models.CharField(max_length=32)
    cwd = models.CharField(max_length=255)
    oldpwd = models.CharField(max_length=255)
    command = models.TextField()