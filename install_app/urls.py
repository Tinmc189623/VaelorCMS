from django.urls import path
from . import views

urlpatterns = [
    path('', views.step_license, name='install_index'),
    path('step/1/', views.step_license, name='install_step1'),
    path('step/2/', views.step_setup, name='install_step2'),
    path('step/3/', views.step_install, name='install_step3'),
]
