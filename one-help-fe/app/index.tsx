import { SafeAreaView, ScrollView } from "react-native";
import Welcome from "@/components/ui/Welcome";
import RecentFunsraises from "@/components/ui/RecentFunsraises";
import RecentEvents from "@/components/ui/RecentEvents";

export default function HomeScreen() {
  return (
    <>
      <ScrollView>
        <SafeAreaView style={{ flex: 1 }}>
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
        </SafeAreaView>
      </ScrollView>
    </>
  );
}
