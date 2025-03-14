import { SafeAreaView, ScrollView } from "react-native";
import Welcome from "@/components/ui/Welcome";
import RecentFunsraises from "@/components/ui/RecentFunsraises";
import RecentEvents from "@/components/ui/RecentEvents";
import RecentRaffles from "@/components/ui/RecentRaffles";

export default function HomeScreen() {
  return (
    <>
      <ScrollView
        bounces={false} // Disables bounce effect on iOS
        overScrollMode="never" // Disables over-scroll effect on Android
        showsVerticalScrollIndicator={false}
      >
        <SafeAreaView className="flex-1 bg-white">
          <Welcome
            user={{
              name: "Олеся",
              donated: 0,
              events: 0,
              raffles: 0,
            }}
          />
          <RecentFunsraises />
          <RecentEvents />
          <RecentRaffles />
        </SafeAreaView>
      </ScrollView>
    </>
  );
}
