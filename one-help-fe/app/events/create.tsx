import React, { useState, useEffect } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  ScrollView,
  Platform,
  KeyboardAvoidingView,
  Alert,
} from "react-native";
import { useRouter } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import { TextInput, Select } from "@/components/controls";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function CreateEvent() {
  const router = useRouter();

  // Provide default values for all fields so they're controlled from start
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      title: "",
      description: "",
      start_date: "",
      end_date: "",
      format: "",
      fundraiseId: "",
      max_participants: "",
      minimum_donation: "",
      address: "",
      imageUrl:""
    },
  });

  const [fundraisers, setFundraisers] = useState([]);

  useEffect(() => {
    const fetchFundraisers = async () => {
      try {
        const token = await AsyncStorage.getItem("token");
        if (!token) return;

        const response = await fetch(
          "http://localhost:8080/api/v0/fundraises/my",
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        const data = await response.json();
        setFundraisers(data.data || []);
      } catch (error) {
        console.error("Error loading fundraisers", error);
      }
    };

    fetchFundraisers();
  }, []);

  const onSubmit = async (data: any) => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const response = await fetch("http://localhost:8080/api/v0/events/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          title: data.title,
          description: data.description,
          startDate: data.start_date
            ? new Date(data.start_date).toISOString()
            : null,
          endDate: data.end_date ? new Date(data.end_date).toISOString() : null,
          format:
            typeof data.format === "string" ? data.format : data.format?.value,
          fundraiseId:
            typeof data.fundraiseId === "string"
              ? data.fundraiseId
              : data.fundraiseId?.id,
          address: data.address,
          maxParticipants: Number(data.max_participants),
          minimumDonation: Number(data.minimum_donation),
          imageUrl: data.imageUrl || "",
        }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Не вдалося створити подію");
      }

      Alert.alert("Успіх", "Подію створено успішно ✅");
      router.push("/events");
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      className="flex-1 bg-white"
    >
      {/* 
        1. Remove TouchableWithoutFeedback 
        2. Use ScrollView with keyboardShouldPersistTaps 
      */}
      <ScrollView
        className="flex-1 p-6"
        keyboardShouldPersistTaps="handled"
        contentContainerStyle={{ paddingBottom: 24 }}
      >
        <Text className="text-2xl font-bold text-gray-msg mb-4">
          Створити подію
        </Text>

        <View className="flex flex-col gap-4">
          {/* Title */}
          <Controller
            control={control}
            name="title"
            rules={{ required: "Назва події є обов'язковою" }}
            render={({ field: { onChange, value } }) => (
              <TextInput
                label="Назва події"
                required="*"
                placeholder="Благодійний марафон"
                placeholderTextColor="#9CA3AF"
                onChangeText={onChange}
                value={value}
                error={!!errors.title}
                errorMessage={errors.title?.message}
              />
            )}
          />

          {/* Description */}
          <Controller
            control={control}
            name="description"
            rules={{ required: "Опис є обов'язковим" }}
            render={({ field: { onChange, value } }) => (
              <TextInput
                label="Опис"
                required="*"
                placeholder="Детально опишіть, що буде на заході"
                placeholderTextColor="#9CA3AF"
                multiline
                onChangeText={onChange}
                value={value}
                error={!!errors.description}
                errorMessage={errors.description?.message}
              />
            )}
          />

          {/* Start Date */}
          <Controller
            control={control}
            name="start_date"
            rules={{ required: "Дата початку є обов'язковою" }}
            render={({ field: { onChange, value } }) => (
              <TextInput
                label="Дата початку"
                required="*"
                placeholder="2025-04-15"
                placeholderTextColor="#9CA3AF"
                onChangeText={onChange}
                value={value}
                error={!!errors.start_date}
                errorMessage={errors.start_date?.message}
              />
            )}
          />

          {/* End Date */}
          <Controller
            control={control}
            name="end_date"
            render={({ field: { onChange, value } }) => (
              <TextInput
                label="Дата завершення"
                placeholder="2025-04-20"
                placeholderTextColor="#9CA3AF"
                onChangeText={onChange}
                value={value}
              />
            )}
          />

          <Controller
            control={control}
            name="format"
            rules={{ required: "Формат події є обов'язковим" }}
            render={({ field: { onChange, value } }) => (
              <Select
                label="Формат події"
                required="*"
                data={[
                  { label: "Онлайн", value: "ONLINE" },
                  { label: "Офлайн", value: "IN_PERSON" },
                  { label: "Гібрид", value: "HYBRID" },
                ]}
                labelField="label"
                valueField="value"
                placeholder="Оберіть формат"
                value={value}
                onChange={(selected: any) => onChange(selected.value)}
                style={{
                  borderWidth: 1,
                  borderColor: "#D1D5DB",
                  paddingHorizontal: 12,
                  paddingVertical: 10,
                  backgroundColor: "white",
                }}
                error={!!errors.format}
                errorMessage={errors.format?.message}
              />
            )}
          />

          {/* Fundraiser Dropdown */}
          <Controller
            control={control}
            name="fundraiseId"
            rules={{ required: "Оберіть збір для події" }}
            render={({ field: { onChange, value } }) => (
              <Select
                label="Збір коштів"
                required="*"
                data={fundraisers}
                labelField="title"
                valueField="id"
                placeholder="Оберіть збір"
                value={value}
                onChange={onChange}
                style={{
                  borderWidth: 1,
                  borderColor: "#D1D5DB",
                  borderRadius: 8,
                  paddingHorizontal: 12,
                  paddingVertical: 10,
                  backgroundColor: "white",
                }}
                error={!!errors.fundraiseId}
                errorMessage={errors.fundraiseId?.message}
              />
            )}
          />

          {/* Max Participants */}
          <Controller
            control={control}
            name="max_participants"
            rules={{ required: "Вкажіть максимальну кількість учасників" }}
            render={({ field: { onChange, value } }) => (
              <TextInput
                label="Максимальна кількість учасників"
                required="*"
                placeholder="100"
                placeholderTextColor="#9CA3AF"
                keyboardType="numeric"
                onChangeText={onChange}
                value={value}
                error={!!errors.max_participants}
                errorMessage={errors.max_participants?.message}
              />
            )}
          />

          {/* Min Donation */}
          <Controller
            control={control}
            name="minimum_donation"
            rules={{ required: "Вкажіть мінімальний внесок" }}
            render={({ field: { onChange, value } }) => (
              <TextInput
                label="Мінімальний внесок (грн)"
                required="*"
                placeholder="50"
                placeholderTextColor="#9CA3AF"
                keyboardType="numeric"
                onChangeText={onChange}
                value={value}
                error={!!errors.minimum_donation}
                errorMessage={errors.minimum_donation?.message}
              />
            )}
          />

          {/* Address */}
          <Controller
            control={control}
            name="address"
            render={({ field: { onChange, value } }) => (
              <TextInput
                label="Адреса (для offline заходів)"
                placeholder="Київ, вул. Хрещатик, 1"
                placeholderTextColor="#9CA3AF"
                onChangeText={onChange}
                value={value}
                error={false}
              />
            )}
          />

<Controller
  control={control}
  name="imageUrl"
  render={({ field: { onChange, value } }) => (
    <TextInput
      label="Зображення (необов'язково)"
      placeholder="https://example.com/image.jpg"
      placeholderTextColor="#9CA3AF"
      onChangeText={onChange}
      value={value}
      error={false}
    />
  )}
/>


        </View>

        <TouchableOpacity
          onPress={handleSubmit(onSubmit)}
          className="bg-primary rounded-lg p-1 mt-6 items-center"
        >
          <Text className="text-white text-lg font-semibold">
            Створити подію 🎉
          </Text>
        </TouchableOpacity>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}
