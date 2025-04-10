import {
  DarkTheme,
  DefaultTheme,
  ThemeProvider,
} from "@react-navigation/native";
import "react-native-reanimated";
import { useFonts } from "expo-font";
import { Stack, useRouter, useSegments } from "expo-router";
import * as SplashScreen from "expo-splash-screen";
import { StatusBar } from "expo-status-bar";
import React, { useEffect } from "react";
import { useColorScheme } from "@/hooks/useColorScheme";
import Navbar from "@/components/Navbar";
import { AuthProvider, useAuth } from "@/context/AuthContext";

import "../global.css";

// Prevent splash screen auto-hide
SplashScreen.preventAutoHideAsync();

function LayoutContent() {
  const colorScheme = useColorScheme();
  const [loaded] = useFonts({
    SpaceMono: require("../assets/fonts/SpaceMono-Regular.ttf"),
  });

  const { isLoggedIn, loading } = useAuth();
  const segments = useSegments();
  const router = useRouter();

  const isAuthPage =
    segments[0] === "auth" &&
    (segments[1] === "register" || segments[1] === "login");

  useEffect(() => {
    if (loaded) {
      SplashScreen.hideAsync();
    }
  }, [loaded]);

  useEffect(() => {
    if (loading) return;

    if (!isLoggedIn && !isAuthPage) {
      router.replace("/auth/login");
    }

    if (isLoggedIn && isAuthPage) {
      router.replace("/");
    }
  }, [isLoggedIn, loading, isAuthPage]);

  if (!loaded || loading) return null;

  return (
    <ThemeProvider value={colorScheme === "dark" ? DarkTheme : DefaultTheme}>
      {!isAuthPage && <Navbar />}

      <Stack>
        <Stack.Screen name="index" options={{ headerShown: false }} />
        <Stack.Screen name="auth/register" options={{ headerShown: false }} />
        <Stack.Screen name="auth/login" options={{ headerShown: false }} />
        <Stack.Screen
          name="fundraises/index"
          options={{ headerShown: false }}
        />
        <Stack.Screen name="fundraises/[id]" options={{ headerShown: false }} />
        <Stack.Screen
          name="fundraises/create"
          options={{ headerShown: false }}
        />
        <Stack.Screen name="events/index" options={{ headerShown: false }} />
        <Stack.Screen name="events/[id]" options={{ headerShown: false }} />
        <Stack.Screen name="events/create" options={{ headerShown: false }} />
        <Stack.Screen name="raffles/index" options={{ headerShown: false }} />
        <Stack.Screen name="raffles/create" options={{ headerShown: false }} />
        <Stack.Screen name="raffles/[id]" options={{ headerShown: false }} />
        <Stack.Screen name="settings/index" options={{ headerShown: false }} />
        <Stack.Screen name="+not-found" />
      </Stack>

      <StatusBar style="auto" />
    </ThemeProvider>
  );
}

export default function RootLayout() {
  return (
    <AuthProvider>
      <LayoutContent />
    </AuthProvider>
  );
}
