import torch
from torch import nn

import lightning as pl
from lightning import Trainer
from lightning.pytorch.callbacks.early_stopping import EarlyStopping

import torch.nn.functional as F
from torchvision import models

from dataset import get_datalaoder

def metric(anchor, positive, negative):
  positive_distance = F.cosine_similarity(anchor, positive)
  negative_distance = F.cosine_similarity(anchor, negative)

  sim_score = torch.sum(positive_distance) / (anchor.size()[0])
  dis_score = torch.sum(negative_distance) / (anchor.size()[0])
  
  delta_score = sim_score - dis_score

  return sim_score, dis_score, delta_score

class TripletModel(pl.LightningModule):
  def __init__(self) -> None:
    super().__init__()
    self.model = models.resnet50(num_classes=128)
    self.criterion = nn.TripletMarginWithDistanceLoss(distance_function=lambda x, y: 1.0 - F.cosine_similarity(x, y))
    self.lr = 1e-5

  def forward(self, image):
    return self.model(image)

  def training_step(self, batch, _):
    a, p, n = batch
    anchor = self.model(a)
    positive = self.model(p)
    negative = self.model(n)

    loss = self.criterion(anchor, positive, negative)
    self.log('train_loss', loss.item(), prog_bar=True, on_epoch=True, on_step=False, sync_dist=True)

    return loss
  
  def validation_step(self, batch, _):
    a, p, n = batch
    anchor = self.model(a)
    positive = self.model(p)
    negative = self.model(n)

    loss = self.criterion(anchor, positive, negative)
    self.log('val_loss', loss.item(), prog_bar=True, on_epoch=True, on_step=False, sync_dist=True)

    sim, dis, delta = metric(anchor, positive, negative)
    metric_similarity = {
      'similar_distance': sim.item(),
      'dissimilar_distance': dis.item(),
      'similarity_delta': delta.item()
    }
    self.log_dict(metric_similarity, on_epoch=True, on_step=False, sync_dist=True)

    return loss

  def configure_optimizers(self):
    return torch.optim.AdamW(self.model.parameters(), lr=self.lr)
    
if __name__ == '__main__':
  torch.set_float32_matmul_precision('high')

  train, test = get_datalaoder(num_workers=6, batch_size=16)
  model = TripletModel()

  trainer = Trainer(max_epochs=100, devices=-1, callbacks=[
    EarlyStopping(monitor="val_loss", mode="min", patience=30)
  ])

  trainer.fit(model, train_dataloaders=train, val_dataloaders=test)