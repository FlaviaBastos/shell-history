import logging

from django.core.paginator import EmptyPage, PageNotAnInteger, Paginator
from django.shortcuts import render
from django.template import loader
from django.http import HttpResponse
from .models import Command


logger = logging.getLogger(__name__)

def index(request):

    search_query = request.GET.get('command', None)
    if search_query:
        latest_commands = Command.objects.filter(command__contains=search_query).order_by('-timestamp')
    else:
        latest_commands = Command.objects.order_by('-timestamp')
    
    paginator = Paginator(latest_commands, 150)

    page = request.GET.get('page')
    commands_paginated = paginator.get_page(page)
    template = loader.get_template('history/commands.html')
    context = {
        'commands_paginated': commands_paginated,
    }
    if search_query:
        context['search_query'] = search_query
    return HttpResponse(template.render(context, request))
