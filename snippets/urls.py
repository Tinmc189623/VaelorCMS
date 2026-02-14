from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name='code_index'),
    path('<int:pk>/', views.detail, name='code_detail'),
    path('submit/', views.submit, name='code_submit'),
]
