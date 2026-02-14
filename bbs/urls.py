from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name='bbs_index'),
    path('<int:pk>/', views.detail, name='bbs_detail'),
    path('post/', views.post, name='bbs_post'),
]
