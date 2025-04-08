import React, { useState } from "react";
import { View, Text, TouchableOpacity, ScrollView, Alert } from "react-native";
import { useRouter } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import { TextInput } from "@/components/controls";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { Platform, TextInput as RNTextInput } from 'react-native';

export default function CreateFundraise() {
  const router = useRouter();
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      title: "",
      description: "",
      goal: "",
      image: "",
      endDate: "",
    },
  });

  const onSubmit = async (data: any) => {

    

    try {
      const token = await AsyncStorage.getItem("token");
      const user = await AsyncStorage.getItem("user");

      if (!token || !user) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const parsedUser = JSON.parse(user);
      const organizerId = parsedUser.id;

      const response = await fetch("http://localhost:8080/api/v0/fundraises/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          title: data.title,
          description: data.description,
          targetAmount: Number(data.goal),
          imageUrl: data.image, 
          endDate: new Date(data.endDate).toISOString(),
          organizerId,
        }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Не вдалося створити збір");
      }

      const fundraise = await response.json();
      Alert.alert("Успіх", "Збір успішно створено!");
      console.log("✅ Created fundraise:", fundraise);
      router.push("/fundraises");
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    }
  };

  return (
    <ScrollView className="flex-1 bg-white p-6">
      <Text className="text-2xl font-bold text-gray-msg mb-4">
        Створити збір
      </Text>
      
      <View className="flex flex-col gap-4">
  <Controller
    control={control}
    name="title"
    rules={{ required: "Назва збору є обов'язковою" }}
    render={({ field: { onChange, value } }) => (
      <TextInput
        label="Назва збору"
        required="*"
        placeholder="Дрон для бригади"
        placeholderTextColor="#9CA3AF"
        onChangeText={onChange}
        value={value}
        error={!!errors.title}
        errorMessage={errors.title?.message?.toString()}
      />
    )}
  />

  <Controller
    control={control}
    name="description"
    rules={{ required: "Опис є обов'язковим" }}
    render={({ field: { onChange, value } }) => (
      <TextInput
        label="Опис"
        placeholder="Детально опишіть, на що збір"
        placeholderTextColor="#9CA3AF"
        required="*"
        multiline
        onChangeText={onChange}
        value={value}
        error={!!errors.description}
        errorMessage={errors.description?.message?.toString()}
      />
    )}
  />

  <Controller
    control={control}
    name="goal"
    rules={{ required: "Необхідна сума є обов'язковою" }}
    render={({ field: { onChange, value } }) => (
      <TextInput
        label="Необхідна сума (грн)"
        required="*"
        placeholder="50000"
        placeholderTextColor="#9CA3AF"
        keyboardType="numeric"
        onChangeText={onChange}
        value={value}
        error={!!errors.goal}
        errorMessage={errors.goal?.message?.toString()}
      />
    )}
  />

  <Controller
    control={control}
    name="image"
    rules={{ required: "Посилання на фото є обов'язковим" }}
    render={({ field: { onChange, value } }) => (
      <TextInput
        label="Посилання на фото"
        required="*"
        placeholder="https://example.com/photo.jpg"
        placeholderTextColor="#9CA3AF"
        onChangeText={onChange}
        value={value}
        error={!!errors.image}
        errorMessage={errors.image?.message?.toString()}
      />
    )}
  />

  <Controller
    control={control}
    name="endDate"
    rules={{ required: "Дата завершення є обов'язковою" }}
    render={({ field: { onChange, value } }) => (
      <View className="flex flex-col">
        <Text className="text-gray-700 mb-1">Дата завершення *</Text>
        <RNTextInput
          value={value}
          onChangeText={onChange}
          placeholder="2025-05-27"
          placeholderTextColor="#9CA3AF"
          keyboardType="default"
          className="border border-gray-300 rounded-lg p-3 bg-gray-100 text-gray-700"
          {...(Platform.OS === 'web' ? { type: 'date' } : {})}
        />
        {errors.endDate && (
          <Text className="text-red-500 text-sm mt-1">
            {errors.endDate.message?.toString()}
          </Text>
        )}
      </View>
    )}
  />
</View>


      <TouchableOpacity
        onPress={handleSubmit(onSubmit)}
        className="bg-primary rounded-lg p-1 mt-6 items-center"
      >
        <Text className="text-white text-lg font-semibold">
          Створити збір ✅
        </Text>
      </TouchableOpacity>
    </ScrollView>
  );
}
