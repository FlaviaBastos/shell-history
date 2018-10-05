from django.db import models
import datetime

# Create your models here.


class Command(models.Model):
    hostname = models.CharField(max_length=255)
    timestamp = models.DateTimeField()
    username = models.CharField(max_length=32)
    altusername = models.CharField(max_length=32, blank=True)
    cwd = models.CharField(max_length=255)
    oldpwd = models.CharField(max_length=255, blank=True)
    command = models.TextField()
    exitcode = models.IntegerField(default=0)
    def __str__(self):
        string = '{user}@{hostname}-{time}'.format(
            user=self.username, hostname=self.hostname,
            time=self.timestamp.strftime('%d%m%Y+%H%M%S:%f'))
        return string
