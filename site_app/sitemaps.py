"""站点地图 - 供搜索引擎抓取"""
from django.contrib.sitemaps import Sitemap
from django.urls import reverse
from articles.models import Article
from bbs.models import BbsPost
from snippets.models import CodeSnippet
from .models import Page


def _reverse(name, *args):
    return reverse(name, args=args if args else ())


class StaticSitemap(Sitemap):
    priority = 0.8
    changefreq = 'daily'

    def items(self):
        return ['home', 'bbs_index', 'code_index', 'articles_index', 'search', 'help', 'faq', 'about', 'games']

    def location(self, item):
        return reverse(item)


class ArticleSitemap(Sitemap):
    changefreq = 'weekly'
    priority = 0.7

    def items(self):
        return Article.objects.filter(status='published').order_by('-updated_at')

    def location(self, obj):
        return _reverse('article_detail', obj.pk)

    def lastmod(self, obj):
        return obj.updated_at or obj.created_at


class BbsSitemap(Sitemap):
    changefreq = 'daily'
    priority = 0.6

    def items(self):
        return BbsPost.objects.filter(approved=True).order_by('-id')[:500]

    def location(self, obj):
        return _reverse('bbs_detail', obj.pk)

    def lastmod(self, obj):
        return obj.created_at


class CodeSitemap(Sitemap):
    changefreq = 'weekly'
    priority = 0.5

    def items(self):
        return CodeSnippet.objects.all().order_by('-id')[:300]

    def location(self, obj):
        return _reverse('code_detail', obj.pk)

    def lastmod(self, obj):
        return obj.created_at


class PageSitemap(Sitemap):
    changefreq = 'monthly'
    priority = 0.6

    def items(self):
        return Page.objects.filter(is_published=True).order_by('order', 'slug')

    def location(self, obj):
        return _reverse('page_detail', obj.slug)

    def lastmod(self, obj):
        return obj.updated_at
