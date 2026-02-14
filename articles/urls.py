from django.urls import path
from . import views
from .feeds import ArticleFeed

urlpatterns = [
    path('feed/', ArticleFeed(), name='article_feed'),
    path('', views.index, name='articles_index'),
    path('<int:pk>/', views.detail, name='article_detail'),
    path('<int:pk>/comment/', views.add_comment, name='article_add_comment'),
    path('<int:pk>/like/', views.like_article, name='article_like'),
    path('<int:pk>/favorite/', views.favorite_article, name='article_favorite'),
    path('new/', views.create, name='article_create'),
    path('<int:pk>/edit/', views.edit, name='article_edit'),
    path('save/', views.save_article, name='article_save'),
]
