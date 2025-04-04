import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { useRouter } from 'expo-router';

export default function RegisterPage() {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [city, setCity] = useState('');
  const [post, setPost] = useState('');
  const [postDepartment, setPostDepartment] = useState('');
  const [website, setWebsite] = useState('');

  const router = useRouter();

  const handleRegister = async () => {
    if (!firstName || !lastName || !phoneNumber || !email || !password) {
      Alert.alert('Помилка', 'Усі обовʼязкові поля мають бути заповнені');
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/api/v0/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          firstName,
          lastName,
          phoneNumber,
          email,
          password,
          city,
          post,
          postDepartment: postDepartment.trim(),
          website,
        }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Помилка під час реєстрації');
      }

      const data = await response.json();
      Alert.alert('Успіх', 'Реєстрація пройшла успішно');
      router.push('/auth/login');
    } catch (error: any) {
      Alert.alert('Помилка', error.message);
    }
  };

  return (
    <ScrollView contentContainerStyle={{ flexGrow: 1, justifyContent: 'center', padding: 40 }} className="bg-white">
      <View className="items-center mb-6">
        <Text className="text-4xl font-bold text-blue-600">💙 OneHelp</Text>
        <Text className="text-2xl font-semibold text-gray-800 mt-2">Реєстрація</Text>
      </View>

      <TextInput placeholder="Ім'я *" value={firstName} onChangeText={setFirstName} className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700" />
      <TextInput placeholder="Прізвище *" value={lastName} onChangeText={setLastName} className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700" />
      <TextInput placeholder="Телефон *" value={phoneNumber} onChangeText={setPhoneNumber} className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700" />
      <TextInput placeholder="Email *" value={email} onChangeText={setEmail} className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700" />
      <TextInput placeholder="Пароль *" value={password} onChangeText={setPassword} secureTextEntry className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700" />

      {/* <TextInput placeholder="Місто" value={city} onChangeText={setCity} className="border border-gray-300 rounded-lg p-3 mt-4 bg-gray-100 text-gray-700" />
      <TextInput placeholder="Пошта (наприклад: NP, UKRPOST)" value={post} onChangeText={setPost} className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700" />
      <TextInput placeholder="Номер відділення (лише цифра)" value={postDepartment} onChangeText={setPostDepartment} keyboardType="numeric" className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700" /> */}
      {/* <TextInput placeholder="Вебсайт (необов'язково)" value={website} onChangeText={setWebsite} className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700" /> */}

      <TouchableOpacity
        onPress={handleRegister}
        className="bg-blue-600 py-1 rounded-xl mt-6 items-center"
      >
        <Text className="text-white text-lg font-semibold">Зареєструватися 🚀</Text>
      </TouchableOpacity>

      <TouchableOpacity onPress={() => router.push('/auth/login')} className="mt-4 self-center">
        <Text className="text-blue-600 text-base font-medium">Вже маєте акаунт? Увійти</Text>
      </TouchableOpacity>

    </ScrollView>
  );
}
