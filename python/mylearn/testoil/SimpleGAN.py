from __future__ import print_function, division
import os
import tensorflow as tf

os.environ["CUDA_VISIBLE_DEVICES"] = "-1"
from keras.datasets import mnist
from keras.layers import Input, Dense, Reshape, Flatten, Dropout
from keras.layers import BatchNormalization, Activation, ZeroPadding2D
from keras.layers.advanced_activations import LeakyReLU
from keras.layers.convolutional import UpSampling2D, Conv2D
from keras.models import Sequential, Model
from tensorflow.keras.optimizers import Adam  # 加上tensorflow
import pandas as pd
import matplotlib.pyplot as plt
import sys
import numpy as np
import csv
import tensorflow.compat.v1 as tf

tf.disable_v2_behavior()


class GAN():
    def __init__(self):
        self.data_rows = 1
        self.data_cols = 2  #
        self.channels = 1
        self.data_shape = (self.data_rows, self.data_cols)
        self.latent_dim = 100
        self.sample_size = 3648
        optimizer = Adam(0.0002, 0.5)
        # 构建和编译判别器
        self.discriminator = self.build_discriminator()
        self.discriminator.compile(loss='binary_crossentropy',
                                   optimizer=optimizer,
                                   metrics=['accuracy'])

        # 构建生成器
        self.generator = self.build_generator()

        # 生成器输入噪音，生成假的图片
        z = Input(shape=(self.latent_dim,))
        data = self.generator(z)  # 生成器生成的图片

        # 为了组合模型，只训练生成器,不训练判别器
        self.discriminator.trainable = False

        # 判别器将生成的图像作为输入并确定有效性
        validity = self.discriminator(data)  # 这个是判别器判断生成器生成图片的结果

        # The combined model  (stacked generator and discriminator)
        # 训练生成器骗过判别器
        self.combined = Model(z, validity)
        self.combined.compile(loss='binary_crossentropy', optimizer=optimizer)

    def build_generator(self):
        model = Sequential()

        model.add(Dense(256, input_dim=self.latent_dim))
        model.add(LeakyReLU(alpha=0.2))
        model.add(BatchNormalization(momentum=0.8))

        model.add(Dense(512))
        model.add(LeakyReLU(alpha=0.2))
        model.add(BatchNormalization(momentum=0.8))

        model.add(Dense(1024))
        model.add(LeakyReLU(alpha=0.2))
        model.add(BatchNormalization(momentum=0.8))

        # np.prod(self.img_shape)=28x28x1
        model.add(Dense(np.prod(self.data_shape), activation='tanh'))
        model.add(Reshape(self.data_shape))
        model.summary()

        noise = Input(shape=(self.latent_dim,))
        data = model(noise)

        # 输入噪音，输出图片
        return Model(noise, data)

    def build_discriminator(self):
        model = Sequential()

        model.add(Flatten(input_shape=self.data_shape))
        model.add(Dense(512))
        model.add(LeakyReLU(alpha=0.2))
        model.add(Dense(256))
        model.add(LeakyReLU(alpha=0.2))
        model.add(Dense(1, activation='sigmoid'))
        model.summary()

        data = Input(shape=self.data_shape)
        validity = model(data)

        return Model(data, validity)

    def train(self, epochs, batch_size=128, sample_interval=500):
        # 加载数据集
        data = pd.read_csv("395-2.txt", header=None)  #
        data = np.array(data.values.tolist()).reshape(3648, 1, -1)  #
        # 将数据进行归一化处理
        # data = data / 8 -1
        # data = np.expand_dims(, axis=3)
        # Adversarial ground truths
        valid = np.ones((batch_size, 1))
        fake = np.zeros((batch_size, 1))
        for epoch in range(epochs):
            # ---------------------
            #  训练判别器
            # ---------------------
            # X_train.shape[0]为数据集的数量，随机生成batch_size个数量的随机数，作为数据的索引
            idx = np.random.randint(0, data.shape[0], batch_size)

            # 从数据集随机挑选batch_size个数据，作为一个批次训练
            x = data[idx]
            # 噪音维度(batch_size,100)
            noise = np.random.normal(0, 1, (batch_size, self.latent_dim))

            # 由生成器根据噪音生成假的图片
            gen_x = self.generator.predict(noise)

            # 训练判别器，判别器希望真实图片，打上标签1，假的图片打上标签0
            d_loss_real = self.discriminator.train_on_batch(x, valid)
            # print(d_loss_real)
            d_loss_fake = self.discriminator.train_on_batch(gen_x, fake)
            # print(d_loss_fake)
            d_loss = 0.5 * np.add(d_loss_real, d_loss_fake)
            # print(d_loss)
            # ---------------------
            #  训练生成器
            # ---------------------
            noise = np.random.normal(0, 1, (batch_size, self.latent_dim))
            # Train the generator (to have the discriminator label samples as valid)
            g_loss = self.combined.train_on_batch(noise, valid)
            # gen_x = (gen_x + 1) * 192
            # 打印loss值
            print("%d [D loss: %f, acc: %.2f%%] [G loss: %f]" % (epoch, d_loss[0], 100 * d_loss[1], g_loss))
            # print("data", gen_x)
            # 每sample_interval个epoch保存一次生成图片
            if epoch % sample_interval == 0:
                self.sample_data(epoch)
                if not os.path.exists("gen_model"):
                    os.makedirs("gen_model")
                self.generator.save_weights("gen_model/G_model%d.hdf5" % epoch, True)
                self.discriminator.save_weights("gen_model/D_model%d.hdf5" % epoch, True)

    def data_write_csv(self, epoch, gen_datas, num):
        if not os.path.exists("gen_data"):
            os.makedirs("gen_data")
        if epoch == 666:
            file_name = "gen_test/test.txt"
        else:
            file_name = "gen_data/%d.txt" % epoch
        print(file_name)
        gen_datas = gen_datas.reshape(num, self.data_cols)
        dt = pd.DataFrame(gen_datas)
        dt.to_csv(file_name, index=0)
        # with open(file_name, "w", encoding="utf-8", newline'') as f:
        #     writer = csv.writer(f)
        #     for data in datas:
        #         writer.writerow(data)
        #     print("保存文件成功，处理结束")

    def sample_data(self, epoch):
        # 重新生成一批噪音，维度为(self.sample_size,100)
        noise = np.random.normal(0, 1, (self.sample_size, self.latent_dim))
        gen_datas = self.generator.predict(noise)
        # 将生成的数据反归一化
        gen_datas = (gen_datas + 1) * 192
        self.data_write_csv(epoch, gen_datas, self.sample_size)

    def test(self, gen_nums=200):
        self.generator.load_weights("gen_model/G_model9000.hdf5", by_name=True)
        self.discriminator.load_weights("gen_model/D_model9000.hdf5", by_name=True)
        noise = np.random.normal(0, 1, (gen_nums, self.latent_dim))
        gen_datas = self.generator.predict(noise)
        # 将生成的数据反归一化
        gen_datas = (gen_datas + 1) * 192
        print(gen_datas)
        if not os.path.exists("gen_test"):
            os.makedirs("gen_test")
        self.data_write_csv(666, gen_datas, gen_nums)


if __name__ == '__main__':
    gan = GAN()
    gan.train(epochs=10000, batch_size=256, sample_interval=500)
    gan.test()
