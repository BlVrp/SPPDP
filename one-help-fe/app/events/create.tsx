import React, { useState } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  ScrollView,
  Platform,
  KeyboardAvoidingView,
  Keyboard,
  TouchableWithoutFeedback,
} from "react-native";
import { useRouter } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import DateTimePicker from "@react-native-community/datetimepicker";
import { TextInput } from "@/components/controls";

export default function CreateEvent() {
  const router = useRouter();
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm();

  const [showStartDatePicker, setShowStartDatePicker] = useState(false);
  const [showEndDatePicker, setShowEndDatePicker] = useState(false);

  const onSubmit = (data: any) => {
    console.log("Нова подія:", data);
    router.back();
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      className="flex-1 bg-white"
    >
      <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
        <ScrollView className="flex-1 p-6">
          <Text className="text-2xl font-bold text-gray-msg mb-4">
            Створити подію
          </Text>
          <View className="flex flex-col gap-4">
            <Controller
              control={control}
              name="title"
              rules={{ required: "Назва події є обов'язковою" }}
              render={({ field: { onChange, value } }) => (
                <TextInput
                  label="Назва події"
                  required="*"
                  placeholder="Благодійний марафон"
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
                  required="*"
                  placeholder="Детально опишіть, що буде на заході"
                  multiline
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.description}
                  errorMessage={errors.description?.message?.toString()}
                />
              )}
            />

            <Text className="text-gray-msg font-medium">
              Дата початку{" "}
              <Text className="text-primary font-medium items-right">* </Text>{" "}
            </Text>
            <Controller
              control={control}
              name="start_date"
              rules={{ required: "Дата початку є обов'язковою" }}
              render={({ field: { onChange, value } }) => (
                <>
                  <TouchableOpacity
                    onPress={() => setShowStartDatePicker(true)}
                    className="border border-gray-300 rounded-lg p-3"
                  >
                    <Text className="text-grey-msg">
                      {value
                        ? value.toLocaleDateString()
                        : "Оберіть дату початку"}
                    </Text>
                  </TouchableOpacity>
                  {showStartDatePicker && (
                    <DateTimePicker
                      value={value || new Date()}
                      mode="date"
                      display="default"
                      onChange={(_, selectedDate) => {
                        setShowStartDatePicker(false);
                        if (selectedDate) onChange(selectedDate);
                      }}
                    />
                  )}
                  {errors.start_date && (
                    <Text className="text-red-500 text-sm">
                      {errors.start_date?.message?.toString()}
                    </Text>
                  )}
                </>
              )}
            />

            <Text className="text-gray-msg font-medium"> Дата завершення </Text>
            <Controller
              control={control}
              name="end_date"
              render={({ field: { onChange, value } }) => (
                <>
                  <TouchableOpacity
                    onPress={() => setShowEndDatePicker(true)}
                    className="border border-gray-300 rounded-lg p-3"
                  >
                    <Text className="text-grey-msg">
                      {value
                        ? value.toLocaleDateString()
                        : "Оберіть дату завершення"}
                    </Text>
                  </TouchableOpacity>
                  {showEndDatePicker && (
                    <DateTimePicker
                      value={value || new Date()}
                      mode="date"
                      display="default"
                      onChange={(_, selectedDate) => {
                        setShowEndDatePicker(false);
                        if (selectedDate) onChange(selectedDate);
                      }}
                    />
                  )}
                </>
              )}
            />

            <Controller
              control={control}
              name="format"
              rules={{ required: "Формат є обов'язковим" }}
              render={({ field: { onChange, value } }) => (
                <TextInput
                  label="Формат події (online/offline)"
                  required="*"
                  placeholder="offline"
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.format}
                  errorMessage={errors.format?.message?.toString()}
                />
              )}
            />

            <Controller
              control={control}
              name="max_participants"
              rules={{ required: "Вкажіть максимальну кількість учасників" }}
              render={({ field: { onChange, value } }) => (
                <TextInput
                  label="Максимальна кількість учасників"
                  required="*"
                  placeholder="100"
                  keyboardType="numeric"
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.max_participants}
                  errorMessage={errors.max_participants?.message?.toString()}
                />
              )}
            />

            <Controller
              control={control}
              name="minimum_donation"
              rules={{ required: "Вкажіть мінімальний внесок" }}
              render={({ field: { onChange, value } }) => (
                <TextInput
                  label="Мінімальний внесок (грн)"
                  required="*"
                  placeholder="50"
                  keyboardType="numeric"
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.minimum_donation}
                  errorMessage={errors.minimum_donation?.message?.toString()}
                />
              )}
            />

            <Controller
              control={control}
              name="address"
              rules={{ required: "Адреса є обов'язковою для offline подій" }}
              render={({ field: { onChange, value } }) => (
                <TextInput
                  label="Адреса (для offline заходів)"
                  required="*"
                  placeholder="Київ, вул. Хрещатик, 1"
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.address}
                  errorMessage={errors.address?.message?.toString()}
                />
              )}
            />
          </View>

          <TouchableOpacity
            onPress={handleSubmit(onSubmit)}
            className="bg-primary rounded-lg p-4 mt-6 items-center"
          >
            <Text className="text-white text-lg font-semibold">
              Створити подію 🎉
            </Text>
          </TouchableOpacity>
        </ScrollView>
      </TouchableWithoutFeedback>
    </KeyboardAvoidingView>
  );
}
