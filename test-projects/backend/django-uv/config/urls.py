"""
URL configuration for django-uv-test project.
"""
from django.http import JsonResponse
from django.urls import path


def hello_world(request):
    return JsonResponse({
        'message': 'Hello from Django with uv!',
        'dependency_manager': 'uv',
        'framework': 'Django'
    })


urlpatterns = [
    path('', hello_world),
    path('health/', hello_world),
]
