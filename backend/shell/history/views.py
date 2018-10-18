from django.core.paginator import EmptyPage, PageNotAnInteger, Paginator
from django.shortcuts import render
from django.template import loader
from django.http import HttpResponse
from .models import Command


def index(request):
    latest_commands = Command.objects.order_by('-timestamp')
    paginator = Paginator(latest_commands, 150)

    page = request.GET.get('page')
    commands_paginated = paginator.get_page(page)
    template = loader.get_template('history/commands.html')
    context = {
        'commands_paginated': commands_paginated,
    }
    return HttpResponse(template.render(context, request))
