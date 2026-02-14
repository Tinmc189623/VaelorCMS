from django.urls import path
from . import api_views

urlpatterns = [
    path('stats/', api_views.stats),
    path('search/', api_views.search),
    path('articles/', api_views.articles_list),
    path('upgrade/', api_views.upgrade),
    path('health/', api_views.health),
]
