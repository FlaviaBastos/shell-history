from django.shortcuts import render
from django.template import loader

# Create your views here.

from django.http import HttpResponse
from .models import Command


def index(request):
    latest_commands = Command.objects.order_by('-timestamp')
    template = loader.get_template('history/commands.html')
    context = {
        'latest_commands': latest_commands,
    }
    return HttpResponse(template.render(context, request))
