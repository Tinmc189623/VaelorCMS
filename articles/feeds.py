"""文章 RSS 订阅"""
from django.contrib.syndication.views import Feed
from django.urls import reverse
from .models import Article


class ArticleFeed(Feed):
    description = '最新发布的文章'

    def title(self, obj=None):
        try:
            from site_app.settings_service import get
            return (get('site_name', '') or 'VaelorCMS').strip() or 'VaelorCMS'
        except Exception:
            return 'VaelorCMS'

    def link(self, obj=None):
        return reverse('articles_index')

    def items(self):
        return Article.objects.filter(status='published').only(
            'id', 'title', 'content', 'created_at'
        ).order_by('-created_at')[:50]

    def item_title(self, item):
        return item.title

    def item_description(self, item):
        content = item.content or ''
        return content[:500] + ('...' if len(content) > 500 else '')

    def item_link(self, item):
        return reverse('article_detail', args=[item.pk])

    def item_pubdate(self, item):
        return item.created_at
