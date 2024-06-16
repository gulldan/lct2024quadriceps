from functools import partial

import torch
from timm.models.layers import trunc_normal_
from timm.models.vision_transformer import VisionTransformer, _cfg
from torch import nn


class VisionTransformerPattern(VisionTransformer):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.pattern_token = nn.Parameter(torch.zeros(1, 1, self.embed_dim))
        num_patches = self.patch_embed.num_patches
        self.pos_embed = nn.Parameter(torch.zeros(1, num_patches + 1, self.embed_dim))
        self.pos_embed_1 = nn.Parameter(torch.zeros(1, 1, self.embed_dim))
        self.head_pat = nn.Linear(self.embed_dim, 90)

        trunc_normal_(self.pattern_token, std=0.02)
        trunc_normal_(self.pos_embed, std=0.02)
        trunc_normal_(self.pos_embed_1, std=0.02)
        self.head_pat.apply(self._init_weights)

    def forward_features(self, x):
        B = x.shape[0]

        cls_tokens = self.cls_token.expand(B, -1, -1)  # stole cls_tokens impl from Phil Wang, thanks
        pattern_token = self.pattern_token.expand(B, -1, -1)
        x = torch.cat((cls_tokens, x, pattern_token), dim=1)

        x = x + torch.cat((self.pos_embed, self.pos_embed_1), dim=1)
        x = self.pos_drop(x)

        for blk in self.blocks:
            x = blk(x)

        x = self.norm(x)
        return x[:, 0], x[:, -1]

    def forward(self, x):
        x = self.patch_embed(x)
        x, x_pat = self.forward_features(x)
        x = self.head(x)
        x_pat = self.head_pat(x_pat)

        return x, x_pat


class VisionTransformer(nn.Module):
    def __init__(
        self,
        pretrained=True,
        cut_at_pooling=False,
        num_features=0,
        norm=False,
        dropout=0,
        num_classes=0,
        dev=None,
    ):
        super(VisionTransformer, self).__init__()
        self.pretrained = True
        self.cut_at_pooling = cut_at_pooling
        self.num_classes = num_classes

        vit = VisionTransformerPattern(
            patch_size=16,
            embed_dim=768,
            depth=12,
            num_heads=12,
            mlp_ratio=4,
            qkv_bias=True,
            norm_layer=partial(nn.LayerNorm, eps=1e-6),
        )
        vit.default_cfg = _cfg()

        vit.head = nn.Sequential()
        self.vit = vit.cuda()
        self.vit.patch_embed.proj = nn.Conv2d(9, 768, kernel_size=(16, 16), stride=(16, 16))
        # self.patch_embed = vit.patch_embed
        # self.forward_features = vit.forward_features
        # self.head_pat = vit.head_pat
        """
        self.base = nn.Sequential(
            vit
        ).cuda()
        """

        self.linear = nn.Linear(vit.embed_dim, 512)

        # self.classifier = build_metric('cos', vit.embed_dim, self.num_classes, s=64, m=0.35).cuda()
        # self.classifier_1 = build_metric('cos', 512, self.num_classes, s=64, m=0.6).cuda()

        self.projector_feat_bn = nn.Sequential(nn.Identity()).cuda()

        self.projector_feat_bn_1 = nn.Sequential(self.linear, nn.Identity()).cuda()

    def forward(self, x):
        x_input = torch.cat((x, x, x), dim=1)
        x_input = self.vit.patch_embed(x_input)

        """
        x = self.vit.patch_embed(x)
        support = self.vit.patch_embed(support)
        support_o = self.vit.patch_embed(support_o)
        x_input = torch.cat((x, support, support_o), dim = 1)
        """

        x, x_pattern = self.vit.forward_features(x_input)
        # x_pattern = self.vit.head_pat(x_pattern)
        x = x.view(x.size(0), -1)
        bn_x = self.projector_feat_bn(x)

        # prob = self.classifier(bn_x, y)
        bn_x_512 = self.projector_feat_bn_1(bn_x)
        # prob_1 = self.classifier_1(bn_x_512, y)

        return bn_x_512, x_pattern
